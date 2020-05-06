A few steps to get started…

Download the cvetracker.zip file

Make sure you have Go installed if you want to run this.  I left the deprecated Python/JSON work in there so you can see some of the evolution, but it is basically noise.

Once you unzip the file and have Go installed go to the ./cvetracker-master/golang/cvetracker directory and execute ./cvetracker

You should see a reply of: expected search, updater or write_example

./cvetracker/write_example

This will create your two default YAML files for storing local configurations.  

Edit the cveupdater.yaml file to include your credentials on the MongoDB.  If you do not have the credentials for the MongoDB server please email douglas.gourlay@gmail.com and I can provide them for you.  You will get two sets of credentials: 
A Read Only credential for populating into cvesearch.yaml
A RW credential for a specific collection within the database i.e. cve_yourname that is ‘your’ collection within the database.

Put the RW credential in the cveupdater.yaml file and be sure to, for the initial collection population set init=true

Then you will run ./cvetracker updater - this takes a long time, about an hour and a half but it builds a transformed and simplified version of the NVD database for network switch devices categorized by vendor.  I suppose this could have been  multi-threaded, spawning a thread for each year or such, but as you really only need to ruin this once it wasn’t worth the optimization.

Once this has been fully run you can then run ./cvetracker search and execute searches against the data set.

You will then modify cveupdater.yaml changing the init Boolean to 'False' and then I would recommend setting a cron job or Lambda function call to run the program daily to keep the database coherent.



