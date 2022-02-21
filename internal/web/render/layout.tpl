<!DOCTYPE html>
<html lang="en">
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,minimum-scale=1,initial-scale=1">
  <title>{{ .Title }}</title>

  {{ .Head }}

  {{- if ne .GTMID "" }}
  <script>
    (function (w, d, s, l, i) {
      w[l] = w[l] || []; w[l].push({ "gtm.start": new Date().getTime(), event: "gtm.js" });
      var f = d.getElementsByTagName(s)[0], j = d.createElement(s), dl = l != "dataLayer" ? "&l=" + l : "";
      j.async = true; j.src = "https://a7s.seankhliao.com/gtm.js?id=" + i + dl;
      f.parentNode.insertBefore(j, f);
    })(window, document, "script", "dataLayer", "{{ .GTMID }}");
  </script>
  {{- end }}

  <link rel="canonical" href="{{ .URLCanonical }}">
  <link rel="manifest" href="/manifest.json">

  <meta name="theme-color" content="#000000">
  <meta name="description" content="{{ .Description }}">

  <link rel="icon" href="https://seankhliao.com/favicon.ico">
  <link rel="icon" href="https://seankhliao.com/static/icon.svg" type="image/svg+xml" sizes="any">
  <link rel="apple-touch-icon" href="https://seankhliao.com/static/icon-192.png">

  <style>
    {{ template "basecss" . }}
    {{ .Style }}
  </style>

  {{- if ne .GTMID "" }}
  <noscript><iframe src="https://a7s.seankhliao.com/ns.html?id={{ .GTMID }}" height="0" width="0" style="display: none; visibility: hidden"></iframe></noscript>
  {{- end }}

  <h1>{{ if .H1 }}{{ .H1 }}{{ else if .Compact }}{{ .Title }}{{ end }}</h1>
  <h2>{{ if .H2 }}{{ .H2 }}{{ else if .Compact }}{{ .Description }}{{ end }}</h2>

  <hgroup>
    {{ if .Compact }}
    <a href="https://seankhliao.com/?utm_medium=sites&utm_source={{.URLCanonical}}">
    {{ else }}
    <a href="/">
    {{ end }}
      <span>S</span><span>E</span><span>A</span><span>N</span>
      <em>K</em><em>.</em><em>H</em><em>.</em>
      <span>L</span><span>I</span><span>A</span><span>O</span>
    </a>
  </hgroup>

  {{ .Main }}

  <footer>
    <a href="https://seankhliao.com/{{ if .Compact }}?utm_medium=sites&utm_source={{.URLCanonical}}{{ end }}">home</a>
    |
    <a href="https://seankhliao.com/blog/{{ if .Compact }}?utm_medium=sites&utm_source={{.URLCanonical}}{{ end }}">blog</a>
    |
    <a href="https://github.com/seankhliao">github</a>
  </footer>
</html>
