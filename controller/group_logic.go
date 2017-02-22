package controller

func GetQRCodeUrlListFromType(t int64) ([]QRCodeUrlInfo, error) {
	tgs := &TypeGroupSetting{
		Type: t,
	}
	has, err := GetTypeGroupSetting(tgs)
	if err != nil {
		plog.Errorf("get type group setting error: %v", err)
		return nil, err
	}
	if !has {
		plog.Errorf("has none this type[%d] of group setting", t)
		return nil, err
	}
	list, err := GetQRCodeUrlInfoList(tgs.GroupNum, t)
	if err != nil {
		plog.Errorf("get qrcode info list error: %v", err)
		return nil, err
	}
	return list, nil
}
