package incident

import (
	"bytes"
	"encoding/json"
	"fmt"
	"incident-tracker/config"
	"incident-tracker/models"
	"incident-tracker/repository"
	"io"
	"net/http"

	"github.com/vladopajic/go-actor/actor"
)

type processor struct {
	application        *config.ApplicationContext
	inC                actor.MailboxReceiver[*models.Incident]
	incidentRepository *repository.IncidentRepository
	requestRepository  *repository.RequestRepository
}

func newProcessor(inC actor.MailboxReceiver[*models.Incident], application *config.ApplicationContext) *processor {
	return &processor{
		application:        application,
		inC:                inC,
		incidentRepository: repository.NewIncidentRepository(application.DB),
		requestRepository:  repository.NewRequestRepository(application.DB),
	}
}

func (w *processor) DoWork(c actor.Context) actor.WorkerStatus {
	select {
	case msg, ok := <-w.inC.ReceiveC():
		if !ok {
			w.application.Logger.Error("channel closed, shutting down")
			return actor.WorkerEnd
		}
		if msg != nil {
			w.application.Logger.Sugar().Infof("Processing incident: %d", msg.ID)
			msg.RequestStatus = models.StatusInProgress
			if err := w.incidentRepository.UpdateIncident(msg); err != nil {
				w.application.Logger.Sugar().Errorf("Error updating incident: %v", err)
				return actor.WorkerContinue
			}

			request, err := GetAIClassification(w.application, msg)
			if err != nil {
				w.application.Logger.Sugar().Errorf("Error getting AI classification: %v", err)
				msg.RequestStatus = models.StatusFailed
				if err := w.incidentRepository.UpdateIncident(msg); err != nil {
					w.application.Logger.Sugar().Errorf("Error updating incident: %v", err)
				}
				return actor.WorkerContinue
			}

			if err := w.requestRepository.CreateRequest(request); err != nil {
				w.application.Logger.Sugar().Errorf("Error creating request: %v", err)
				return actor.WorkerContinue
			}

			var aiResult map[string]string
			if err := json.Unmarshal([]byte(request.ResponseBody), &aiResult); err != nil {
				w.application.Logger.Sugar().Errorf("Error unmarshalling AI result: %v", err)
				return actor.WorkerContinue
			}

			msg.AISeverity = aiResult["ai_severity"]
			msg.AICategory = aiResult["ai_category"]
			msg.RequestStatus = models.StatusCompleted
			if err := w.incidentRepository.UpdateIncident(msg); err != nil {
				w.application.Logger.Sugar().Errorf("Error updating incident: %v", err)
			}
		}

		return actor.WorkerContinue

	case <-c.Done():
		w.application.Logger.Info("is shutting down")
		return actor.WorkerEnd

	}
}

func GetAIClassification(appCtx *config.ApplicationContext, incident *models.Incident) (*models.Request, error) {
	url := appCtx.Config.OpenAI.APIUrl
	apiKey := appCtx.Config.OpenAI.APIKey
	prompt := fmt.Sprintf(`
Given this IT incident:
Title: %s
Description: %s
Affected Service: %s

Classify it into:
- ai_severity: One of [Low, Medium, High, Critical]
- ai_category: One of [Network, Software, Hardware, Security]

Return a JSON like:
{
  "ai_severity": "High",
  "ai_category": "Network"
}
`, incident.Title, incident.Description, incident.AffectedService)

	requestBody := map[string]interface{}{
		"model": appCtx.Config.OpenAI.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": appCtx.Config.OpenAI.Temperature,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := appCtx.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return &models.Request{
			IncidentID:     incident.ID,
			RequestBody:    prompt,
			ResponseBody:   string(body),
			ResponseStatus: uint(resp.StatusCode),
		}, nil
	}
	var openaiResp models.OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, err
	}

	if len(openaiResp.Choices) == 0 {
		return &models.Request{
			IncidentID:     incident.ID,
			RequestBody:    prompt,
			ResponseBody:   "{}",
			ResponseStatus: uint(resp.StatusCode),
		}, nil
	}
	// Parse the JSON returned in message content
	var aiResult map[string]string
	err = json.Unmarshal([]byte(openaiResp.Choices[0].Message.Content), &aiResult)
	if err != nil {
		return nil, err
	}
	// os.WriteFile("abc.json", []byte(openaiResp.Choices[0].Message.Content), 777)
	b, err := json.Marshal(aiResult)
	if err != nil {
		return nil, err
	}
	return &models.Request{
		IncidentID:     incident.ID,
		RequestBody:    prompt,
		ResponseBody:   string(b),
		ResponseStatus: uint(resp.StatusCode),
	}, nil
}
