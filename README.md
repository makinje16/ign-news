# **Overview**
This go service through the use of Twilio and NewsAPI will hit the NewsAPI endpoint for igns headlines and send them to your cellphone number along with a small description of the headline and a link to the article or video.

# **Setup**
For this program to work you need to export these things in your environment:

```
export TWILIO_API_KEY=[api-key]
export TWILIO_SID=[account-sid]
export TWILIO_NUMBER=[your twilio number]
export NEWS_API_KEY=[api-key]
```

Both your Twilio API key and account SID can be found on your twilio console

News API key can be found after signing up with https://newsapi.org/

# **Usage**
The contract for this program is as follows:

```
./news <number to send messages to>
```

Note: the number needs to have the country Dialing code as well.
So for the United States an example number would be 
`+15555555555`