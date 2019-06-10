package models

func CreatePartitionTable(creHashTable, creTableTriggerFunc, creTableTrigger string) error {
	tx := DB.Self.Begin()
	if err := tx.Exec(creHashTable).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec(creTableTriggerFunc).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Exec(creTableTrigger).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

