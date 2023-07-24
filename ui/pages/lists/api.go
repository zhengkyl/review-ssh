package lists

import "github.com/zhengkyl/review-ssh/ui/common"

// TODO use pagination, but for now 50 is more than enough
const reviewsEndpoint = common.ReviewBase + "/reviews?category=Film&per_page=50"
