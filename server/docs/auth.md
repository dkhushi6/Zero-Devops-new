# Current Implementation of Authoriation

### User
I have maintained a seperate contract under domain named as user.go which consists of all the user profile information
It consists of :
ID
Github ID
Username 
Email (Optional if the org is private)
AvatarURL
CreatedAT

And then storing the information in the Postgres Database 

### Auth
Here it handles the two step process into one step where the github oauth is managed along with the callback to installtion id to install the github app for a user

Also we ensure to generate two session tokens access token and refresh token and do no rely on the github app access token for the user session management 

Also a logout feature to allow the user to log out from the registred session 

### Github

Here we are storing the follwing infromation when the github app is installeed and after installtion there is the following attributes we are taking care of : 
ID			
UserID (for a User type)	
InstalltionID	
AccountName	

And then storing the information in the Postgres Database 