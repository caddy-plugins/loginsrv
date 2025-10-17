go install github.com/daaku/cssdalek@latest
baseDir=../template/
cssdalek \
  --css "_bootstrap.min.css"\
  --word "${baseDir}*.html"\
  --include-class "btn-.*" > bootstrap.min.css

cssdalek \
  --css "_font-awesome.css"\
  --word "${baseDir}*.html"\
  --include-class "fa-(github|gitlab|yahoo|wechat|gitea|google|bitbucket|paypal|stripe|salesforce)" > font-awesome.min.css
