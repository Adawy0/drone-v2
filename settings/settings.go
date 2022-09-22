package constants

const MAX_WEIGHT int = 500

func GetDroneModels() map[string]string {
	return map[string]string{"lightweight": "Lightweight", "middleweight": "Middleweight", "cruiserweight": "Cruiserweight", "heavyweight": "Heavyweight"}
}

func GetDroneState() map[string]string {
	return map[string]string{"idle": "IDLE", "loading": "LOADING", "loaded": "LOADED", "delivering": "DELIVERING", "delivered": "DELIVERED", "returning": "RETURNING"}
}
