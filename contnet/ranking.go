package contnet

import "time"

const (
    __GRAVITY = 0.6
    __QUALITY_RANGE_EFFECT = 0.6
    __POPULARITY_EFFECT_MINUTES_PER_POINT = time.Minute * 1
)

func __age(referenceTime time.Time, content Content) time.Time {
    // get real age as duration from content creation to reference time
    realAge := referenceTime.Sub(content.CreatedAt)

    // calculate effect the popularity & quality have on the age
    popularityEffect := __POPULARITY_EFFECT_MINUTES_PER_POINT * time.Duration(content.Popularity*(1+content.Quality*__QUALITY_RANGE_EFFECT))

    // calculate effect the gravity has on the age
    gravityEffect := time.Duration(realAge.Seconds() * __GRAVITY) * time.Second

    // final result: reduce age based on popularity effect but increase it based on gravity effect
	return content.CreatedAt.Add(popularityEffect - gravityEffect)
}
