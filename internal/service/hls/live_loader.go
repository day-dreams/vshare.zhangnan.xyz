package hls

import (
	"context"
	"fmt"
	"github.com/day-dreams/vshare.zhangnan.xyz/internal/service"
	"github.com/day-dreams/vshare.zhangnan.xyz/internal/utils"
	"os/exec"
)

// LiveFfmpegLoader 实时调用ffmpeg，生成m3u8文件
type LiveFfmpegLoader struct {
}

func (l *LiveFfmpegLoader) M3u8Segment(ctx context.Context, req *ReqM3u8Segment) (*ResM3u8Segment, error) {

	vInfo, err := service.VideoInfoGet(ctx, GetVidPath(req.Vid))
	if err != nil {
		return nil, err
	}
	startAt := req.Index * GetDPB()
	resolution := 1080 // todo 动态码率？
	pipe := "pipe:out%03d.ts"
	args := []string{
		"-timelimit", "45",
		"-ss", fmt.Sprintf("%v.00", startAt),
		"-i", vInfo.Path, // The source file
		"-t", fmt.Sprintf("%v.00", GetDPB()), // The duration
		"-vf", fmt.Sprintf("scale=-2:%v", resolution),
		"-vcodec", "libx264",
		"-preset", "veryfast",
		"-acodec", "aac",
		"-pix_fmt", "yuv420p",
		"-force_key_frames", "expr:gte(t,n_forced*5.000)",
		"-f", "ssegment",
		"-segment_time", fmt.Sprintf("%v.00", GetDPB()),
		"-initial_offset", fmt.Sprintf("%v.00", startAt),
		pipe,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	// utils.Logger().Debug(cmd.String())

	segment, err := cmd.Output()
	if err != nil {
		utils.Logger().Errorf("cmd.output failed. %v", err)
		return nil, err
	}
	return &ResM3u8Segment{Content: segment}, nil
}

func (l *LiveFfmpegLoader) M3u8PlayList(ctx context.Context, req *ReqM3u8PlayList) (*ResM3u8PlayList, error) {
	path := GetVidPath(req.Vid)
	vInfo, err := service.VideoInfoGet(ctx, path)
	if err != nil {
		return nil, err
	}

	res := &ResM3u8PlayList{PlayListContent: ""}
	total := vInfo.Duration
	DPB := GetDPB()
	cur := 0

	// header + play list
	res.PlayListContent += fmt.Sprintf("#EXTM3U\n#EXT-X-TARGETDURATION:%d\n\n", DPB)
	for ; float64(cur*DPB) <= total; cur += 1 {
		duration := float64(DPB)
		if float64(cur*DPB)+duration > total {
			duration = total - float64(cur*DPB)
		}
		// [cur*GetDPB,cur*GetDPB+duration)
		res.PlayListContent += fmt.Sprintf("#EXTINF:%.2f,\n%s?vid=%s&segment=%d\n", duration, req.Path, req.Vid, cur)
	}
	res.PlayListContent += fmt.Sprintf("\n#EXT-X-ENDLIST")
	return res, nil
}
