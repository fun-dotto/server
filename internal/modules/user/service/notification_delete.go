package service

func (s *NotificationService) DeleteNotification(id string) error {
	return s.repo.DeleteNotification(id)
}
