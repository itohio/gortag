# gortag
Real Time Audio Generator implemented in Go using Fyne.

Any audio system analysis cannot be performed without a signal generator, the audio system itself and an audio analyzer.
This application works as a functional generator that allows users to generate various test signals.

In order to support the development of audio analysis software, I haven't yet releazed the dependency `gocoustics` that implements all the necessary bits and pieces.
Also, I have split this app into two versions - Free version and Paid version. 
These versions will be available on Google Play and Microsoft App Store.
However, I am also planning to launch a web version of this app.
Of course, I believe, I'll make all builds available as downloadable artifacts in the Releases section.

At some point in time I'll make the `gocoustics` library available too after I have bought a few beers with the revenue from the paid versions or the sponsorship(whichever comes first!)

# Free version

Planned Features:
- [ ] Select an audio output 
- [ ] Select a wav file as an output
- [ ] Support multichannel output (each tone can be played on any combination of channels)
- [ ] Generate a single tone
- [ ] Generate multiple tones at the same time
- [ ] Modulate one tone using another tone (AM, FM)
- [ ] Be able to chain tones that modulate each other
- [ ] Send encoded message
  - [ ] FSK encoding
  - [ ] MFSK encoding

# Paid version
All the same features, except that one can also generate coded test signals (so that the paid version of gortaa can understand the test signal and do appropriate analysis).

- [ ] Be able to generate a coded test signal
- [ ] Be able to build a coded test signal configuration
- [ ] Be able to save/load the tone generation configuration


## Coded test signal
Coded test signal is such a test signal that carries within itself the type and configuration of test tones and purpose of the test. 
For example, one might want to test audio system amplitude non linearities and THD at specific frequencies.
Or one might want to measure the whole 5.1 or 7.1 theater system impulse/frequency complete with delays for each channel separately.
The coded test signal can encode all that so that the goRTAA too on the listening end can properly interpret received test signals and
generate a proper report about the audio system as well as room acoustics properties.

You can imagine it be like a very popular fully automated REW(Room EQ Wizard) tool. 
Completely open source and multiplatform!
