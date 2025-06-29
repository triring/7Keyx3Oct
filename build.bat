mkdir uf2

tinygo build -o uf2/AnalogInput_light.uf2           -target=pico -size short examples/AnalogInput_light/main.go
tinygo build -o uf2/AnalogInput_light_waveform.uf2  -target=pico -size short examples/AnalogInput_light_waveform/main.go
tinygo build -o uf2/hand-cranked_Motor.uf2          -target=pico -size short examples/hand-cranked_Motor/main.go
tinygo build -o uf2/InterruptTest.uf2               -target=pico -size short examples/InterruptTest/main.go
tinygo build -o uf2/L-chika.uf2                     -target=pico -size short examples/L-chika/main.go
tinygo build -o uf2/L-chikaToBeep.uf2               -target=pico -size short examples/L-chikaToBeep/main.go
tinygo build -o uf2/PWM_Motor_Music.uf2             -target=pico -size short examples/PWM_Motor_Music/main.go
tinygo build -o uf2/PWM_MotorCTL.uf2                -target=pico -size short examples/PWM_MotorCTL/main.go
tinygo build -o uf2/SwitchToConsole.uf2             -target=pico -size short examples/SwitchToConsole/main.go
tinygo build -o uf2/SwitchToLED.uf2                 -target=pico -size short examples/SwitchToLED/main.go
tinygo build -o uf2/ToneTest.uf2                    -target=pico -size short examples/ToneTest/main.go

