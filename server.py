import os
import jinja2
import webapp2

from google.appengine.api import channel
from google.appengine.api import users
from google.appengine.ext import ndb

JINJA_ENVIRONMENT = jinja2.Environment(
    loader=jinja2.FileSystemLoader(os.path.dirname(__file__)),
    extensions=['jinja2.ext.autoescape'],
    autoescape=True)

class ActiveUser(ndb.Model):
  userid = ndb.StringProperty()

class ChatHandler(webapp2.RequestHandler):
  def post(self):
    user = users.get_current_user()
    message = self.request.body
    if message:
      formatted_message = '%s: %s' % (user.nickname(), message)
      for activeuser in ActiveUser.query():
        channel.send_message(activeuser.userid, formatted_message)

class DisconnectHandler(webapp2.RequestHandler):
  def post(self):
    userid = self.request.get('from')
    if userid:
      ndb.Key(ActiveUser, userid).delete()

class MainHandler(webapp2.RequestHandler):
  def get(self):
    # Check if user is logged in
    user = users.get_current_user()
    if user:
      # Create unique chat channel for user and save to active list
      token = channel.create_channel(user.user_id())
      # Add user to active users
      ActiveUser(userid=user.user_id(), id=user.user_id()).put()
      # Render index.html for user
      template = JINJA_ENVIRONMENT.get_template('index.html')
      self.response.write(template.render({'token': token}))
    else:
      self.redirect(users.create_login_url(self.request.uri))

app = webapp2.WSGIApplication([
    ('/', MainHandler),
    ('/chat', ChatHandler),
    ('/_ah/channel/disconnected/', DisconnectHandler),
], debug=True)
