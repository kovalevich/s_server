require 'webrick'
require 'json'

class JSONHandler < WEBrick::HTTPServlet::AbstractServlet
  def do_POST(request, response)
    # Set response headers
    response.status = 200
    response['Content-Type'] = 'text/plain'

    # Respond with OK message
    response.body = 'Data received and saved successfully.'
  end
end

server = WEBrick::HTTPServer.new(Port: 443)
server.mount '/', JSONHandler

# Trap interrupt signal (Ctrl-C) to gracefully shut down the server
trap('INT') { server.shutdown }

# Start the server
server.start
