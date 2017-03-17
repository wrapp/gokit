task :test do
  sh "glide -q install"
  sh "go test -v ./kit"
  sh "go test -v ./util"
end

task :run  do
  sh "SERVICE_NAME='example' go run example/example.go"
end
