Favicon
-------
Just call
convert comic_200_200.png  -bordercolor white -border 0 \( -clone 0 -resize 16x16 \) \( -clone 0 -resize 32x32 \) \( -clone 0 -resize 48x48 \) \( -clone 0 -resize 64x64 \) -delete 0 -alpha off -colors 256 favicon.ico
to create the fav icon.

See
https://unix.stackexchange.com/questions/89275/how-to-create-ico-file-with-more-then-one-image
http://www.imagemagick.org/Usage/thumbnails/#favicon
