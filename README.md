# pwrswtch

pwrswtch is a small GTK3 application that connects to a [Hamlib](https://github.com/Hamlib/Hamlib) `rigctld` server and sets the power level and the PTT of the connected transceiver. This is handy if the power setting of your transceiver is hidden within an "elaborate" menu structure and you don't want to tune your antenna with 100 Watts.

## Configuration

pwrswtch loads its configuration from the file `~/.config/hamradio/conf.json` using [ftl/hamradio/cfg](https://github.com/ftl/hamradio). The following JSON snippet contains the default configuration. Insert it into `conf.json` and modify it to your needs.

```json
{
    "pwrswtch": {
        "levels": [
            { "watts": "10", "value": "0.039216" },
            { "watts": "30", "value": "0.117647" },           
            { "watts": "50", "value": "0.196078" },
            { "watts": "100", "value": "1.0" }
        ],
        "tuningValue": "0.039216",
        "trxHost": "localhost:4532"
    }
}
```

`levels` contains a list of power levels. pwrswtch shows one button for each of these levels. The `value` property contains the value that needs to be send via hamlib to the transceiver in order to set the power level to the desired `watts` value. The values in the default configuration are valid for my FT-450D. To find out the right values for your transceiver, you can do the following: 

1. connect to your hamlib server in the terminal (e.g. running on localhost:4532): `netcat localhost 4532`
2. set the power level of your transceiver to the desired amount of Watts
3. query the `value`: `\get_level RFPOWER`
4. copy the returned value into your configuration

`tuningValue` is the power level value that should be used for tuning.

`trxHost` can be used to connect to a hamlib server running on a different host.

## Disclaimer
I develop this tool for myself and just for fun in my free time. If you find it useful, I'm happy to hear about that. If you have trouble using it, you have all the source code to fix the problem yourself (although pull requests are welcome).

## License
This tool is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/) 2019