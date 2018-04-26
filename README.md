Tested devices
-----

Console utility to control Uniel devices over RS-485 interface

 - [UCH-M121RX](https://uniel.ru/catalog/2200)
 - [UCH-M131RC](https://uniel.ru/catalog/2198)

Raspberry Pi zero w
-------------------

GOOS=linux;GOARCH=arm;GOARM=6

Usage
-----

Enable 4 channel on device with address 1
```
UNIEL_DEVICE_ADDRESS=1 uniel_rpi on -c 4
```

Disable 4 channel on device with address 1
```
UNIEL_DEVICE_ADDRESS=1 uniel_rpi off -c 4
```

TODO
----

 - Move device address from `env` to `--device` option
 - Add text aliases for device address and channels (`UCH-M131RC` instead of 1, `kitchen_spots2` instead of 7)
 - Add commands `set` (`dim`), `status`, `setAddress`

Check device status (show each channel mode)
```bash
./uniel --device=UCH-M131RC status
```

Enable channel 2 (set 100% for dimmable devices)
```bash
./uniel --device=UCH-M121RX on 2
./uniel --device=UCH-M131RC set 2 255
```

Set 50% dim on channel 1 for dimmable device
```bash
./uniel --device=UCH-M131RC set 1 127
```

Disable channel 3 (bathroom_ceiling)
```bash
./uniel --device=UCH-M121RX off 3
./uniel --device=UCH-M131RC set 3 0
./uniel --device=UCH-M131RC dim bathroom_ceiling 0
```

