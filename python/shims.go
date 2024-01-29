package python

var pythonShimCall = "$(v python which --raw) $@"

var pipShimCall = "$(v python which --raw) -m pip $@"

var Shims = map[string]string{
	"python":  pythonShimCall,
	"python3": pythonShimCall,
	"pip":     pipShimCall,
	"pip3":    pipShimCall,
}
