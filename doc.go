/*
   opmlfeed is a simple web service that accepts incoming
   OPML feeds for users and provides an easy URL for sharing
   those OPML feeds.

   Users should POST their feed with the file under the 'feed'
   parameter and the uuid under the 'id' parameter. opmlfeed
   will return the url for the user's feed. Users can then
   perform a GET request on that url to retrieve the OPML feed.
*/
/*
   This package is provided under a dual ISC / public domain license.
   The public domain license is the one applicable to the user of this
   code; you are free to choose whichever affords the maximum freedom
   to you.

   --------------------------------------------------------------------
   The ISC license:

   Copyright (c) 2012 Kyle Isom <kyle@tyrfingr.is>

   Permission to use, copy, modify, and distribute this software for any
   purpose with or without fee is hereby granted, provided that the 
   above copyright notice and this permission notice appear in all 
   copies.

   THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL 
   WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED 
   WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE 
   AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL
   DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA
   OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER 
   TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR 
   PERFORMANCE OF THIS SOFTWARE.
   --------------------------------------------------------------------
*/
package main
