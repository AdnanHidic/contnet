package contnet

import (
	"log"
	"testing"
)

var exampleTitle = `Halfway through my 12 weeks of forced FI`

var exampleDescription = `About 6 weeks ago, I was fired without cause. My contract required them to keep me on the payroll until last week.

I was just offered my dream job at a large mutual fund company and my start date is in 6-weeks.

I thought some people would get a kick out of this thread since the question of "what would I do with all my free time?" comes up.

Quick bio:
* Late 20s male * Major City
* Investment professional (not techie!)

You'll be so busy that you won't know how you used to have time for work

I used to sleep 6-7 a night and work 10-12 hours a day. Now I sleep 8-9 hours a night, and I feel amazing. So many people have told me I look younger. I love sleep and not waking to an alarm clock.

You'll also spend more time doing what you already do. For me, I started going to the gym more often and stayed a little longer. I also started cooking - which is something I talked about but avoided for the last 10-years. Finally, I can enjoy the simple things for longer, such as coffee, reading (even reddit), and spending time with friends and family.

The last point is worth reiterating. Being a type A guy, I always felt like I had to budget my time: "okay, we'll drive to the burbs to see the parents, but we need to leave by 2:35pm". Now, it's so much easier to go through the day and go with the flow.

Unique to my situation, I spent some time worrying about my next job and a lot of time rehabbing an investment property. So, as much as I enjoyed the time off, I wasn't drinking margaritas on the beach.

But, now that I'm done with the rehab and I signed my offer letter, I'm going to travel for 6 weeks. First, Hawaii, then Denver, then London, and hopefully 1-2 more places I haven't thought of. M

All in, this reiterated what FI means to me. Unlike many on here who want to retire ASAP, the job offer reminded me how much I love being an investor. I don't want to be jobless - I just want "fuck you" money. However, I learned that who I am is more than what my job is. I could be happy without one.

Lastly, (and always controversially), I don't want to reach FI by going the route of cost cutting. I couldn't live like MMM. My dream retirement is in the lap of luxury ($10k-$15k per month). I'm fortunate that I have a high income job and no desire to retire before 40. Maybe ever (like Buffett).
`

var exampleComments = []string{
	`Even if one doesn't go the FI route, it's probably healthy to take some time off from work occasionally. I mean more than the 2-4 weeks a year that most people do. It allows you to spend some time as yourself rather than "Mr. So-and-so, B.A., M.S. Junior Engineering Project Manager". It would be nice if we could find a way to work periodic sabbaticals into our business culture. `,
	`I have heard of certain firms giving sabbaticals every 5th year for 3 months. I hope to push for that at my firm. I am currently working on getting standing desks. Time off for senior (10+ years) employees is my next goal. Time away can have multiple positive benefits. Besides the RnR, just the mental break can allow you to come back to the office with new vigor and new insight to old problems. It also solidifies that the company cares about its employees' health.`,
	`Up until recently, I thought I wanted to RE at 40 but now I'm thinking I'd really like to just have more time off. I could see myself working the rest of my life if I was taking long vacations and had more PTO. Investing my time in Hobbies, family and friends is extremely important but I do enjoy the challenges/fulfillment that work life has to offer. Just not the 40+ hours a week that it currently requires.`,
	`Still getting into the groove of early retirement, regarding sleep I've found for me the additional shuteye comes from naps instead of at night. For last five years of career (software dev) I worked from home and had pretty flexible hours so didn't use an alarm clock anyway. I guess I hang around in bed more sometimes in the morning but getting up relatively early is just something I do naturally. I'm in my mid-40s so maybe it is an old man thing, I can probably make a spreadsheet and track median scrotum hang length to determine if old man characteristics are happening. Hell I wouldn't mind the gift of crazy old man strength, using flabby wan freely swinging bicep muscles to easily choke out some whipper snapper at bar while shouting angrily. But I nap a lot more now, like after lunch just go lay on the couch and doze or sometimes late afternoon siesta. I like it.`,
	`You're the man. I'm jelly.`,
}

func TestKeywordExtraction(t *testing.T) {
	input := &ContentKeywordExtractionInput{
		Title:       exampleTitle,
		Description: exampleDescription,
		Comments:    exampleComments,
	}
	keywords := ContentKeywordExtractor.Extract(input, 5)

	for i := 0; i < len(keywords); i++ {
		log.Println(keywords[i])
	}
}
