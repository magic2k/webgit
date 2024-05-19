Simple git utility for pulling, branching and rollback repo via web

How to install and run:

   1. git clone https://github.com/magic2k/webgit.git /opt/webgit
   2. cp /opt/webgit/webgit.service /etc/systemd/system/
   3. systemctl daemon-reload
   4. systemctl start webgit

    You should see it's running on http://hostname:8080
