package consumer

import (
	"log"

	"github.com/IBM/sarama"
	cdcservices "github.com/niflheimdevs/backend/internal/application/cdc/services"
)

type KafkaCdc struct {
	TeamCdc    cdcservices.TeamCdc
	ProjectCdc cdcservices.ProjectCdc
	UserCdc    cdcservices.UserCdc
}

func NewKafkaCdc(
	teamCdc cdcservices.TeamCdc,
	projectCdc cdcservices.ProjectCdc,
	userCdc cdcservices.UserCdc,

) *KafkaCdc {
	return &KafkaCdc{
		TeamCdc:    teamCdc,
		ProjectCdc: projectCdc,
		UserCdc:    userCdc,
	}
}

func (kc *KafkaCdc) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group setup")
	return nil
}

func (kc *KafkaCdc) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group cleanup")
	return nil
}

func (kc *KafkaCdc) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		log.Printf(
			"Topic=%s Partition=%d Offset=%d Key=%s Value=%s",
			msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value,
		)

		var err error
		switch msg.Topic {
		case "postgres.public.users":
			err = kc.UserCdc.UserCapturer(msg.Value)
		case "postgres.public.team":
			err = kc.TeamCdc.TeamCapturer(msg.Value)
		case "postgres.public.users_career_tag":
			err = kc.UserCdc.UserTagCapturer(msg.Value)
		case "postgres.public.project_tag":
			err = kc.ProjectCdc.ProjectTagCapturer(msg.Value)
		case "postgres.public.project":
			err = kc.ProjectCdc.ProjectCapturer(msg.Value)
		default:
			log.Printf("what the actual fuck? %s wasn't suppose to be here", msg.Topic)
		}
		if err != nil {
			log.Println("something went wrong. info: ", err)
		}
		sess.MarkMessage(msg, "")

	}
	return nil
}
