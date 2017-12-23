Tested devices
-----

 - UCH-M121RX
 - UCH-M131RC
 
Usage
-----

Check device status
```bash
./uniel --device=UCH-M131RC --channel=1 status
```

Set channel mode to maximum
```bash
./uniel --device=UCH-M131RC --channel=1 on 1
./uniel --device=UCH-M131RC --channel=1 set 1 255
```

Disable channel 
```bash
./uniel --device=UCH-M131RC --channel=1 off 1
./uniel --device=UCH-M131RC --channel=1 set 1 0
```

