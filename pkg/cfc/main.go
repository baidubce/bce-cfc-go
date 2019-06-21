package cfc

func Main() error {
	config, err := NewRuntimeConfig()
	if err != nil {
		return err
	}

	client := NewCfcClient(config, 0)
	err = client.WaitInvoke()
	client.Close()
	if err != nil {
		return err
	}
	return nil
}
