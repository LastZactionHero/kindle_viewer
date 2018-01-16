# Kindle Desktop Viewer

This is an experimental project to use the Kindle web browser as a desktop monitor. If you come across something that would be easier to view on an e-ink display, just read it on your Kindle instead.

The host application periodically takes a screenshot of the desktop, which is displayed in the Kindle browser. If the image has changed, image on the Kindle is updated.

Because of the lower resolution, the Kindle only shows the top corner of the display (700px x 600px).

## To Use

Currently runs on OS X. If you're using a retina display, you may need to scale down the image.

1. `go run server.go`
2. Connect the Kindle to the same network as your computer.
3. In the Kindle browser, go to address `<ip address>:4000`

## TODO

- Update to only modify the individual, changed pixels, to prevent the Kindle full-page refresh.