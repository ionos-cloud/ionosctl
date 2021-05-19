package commands

//
//func TestRunIpFailoverList(t *testing.T) {
//	var b bytes.Buffer
//	w := bufio.NewWriter(&b)
//	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
//		viper.Reset()
//		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
//		viper.Set(config.ArgQuiet, false)
//		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
//		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resources.IpBlock{IpBlock: testIpBlock}, nil, nil)
//		err := RunIpFailoverList(cfg)
//		assert.NoError(t, err)
//	})
//}
//
//func TestRunIpFailoverListErr(t *testing.T) {
//	var b bytes.Buffer
//	w := bufio.NewWriter(&b)
//	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
//		viper.Reset()
//		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
//		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resources.IpBlock{IpBlock: testIpBlock}, nil, testIpBlockErr)
//		err := RunIpFailoverList(cfg)
//		assert.Error(t, err)
//	})
//}
