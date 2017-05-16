# goillotine
Make podcasts headless again

A lot of interesting youtube self-broadcasters put a lot of emphasize on their talking head. I don't want to call out names, but you know who you are. If your video is just you talking about an interesting subject, then the amount of content is no better than a pure audio version that would be much easier to consume for example on a commute.

So this goillotine was intended to make podcasts headless. To consume the youtube video and provide the audio-only version in a format that allows easy subscriptions by your favorite mobile podcast player.

While this project is actually working as intended (in a proof-of-concept kind of way) and the idea is still relevant, I consider this repository a failed attempt. **DON'T USE IT!**

The reason is that the go ecosystem is not up for this task compared to for example python. The packages that do youtube downloading are lacking behind the top dog [youtube-dl](https://github.com/rg3/youtube-dl) and the existing ffmpeg bindings are not as stable as the python-counterpart.

Conclusion: Redo in python, go is currently not the right language for this. Allows for much nicer names though.
