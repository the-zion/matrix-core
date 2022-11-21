// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package biz

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson899f4d6bDecodeGithubComTheZionMatrixCoreAppAchievementServiceInternalBiz(in *jlexer.Lexer, out *CommentMap) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "uuid":
			out.Uuid = string(in.String())
		case "medal":
			out.Medal = string(in.String())
		case "mode":
			out.Mode = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson899f4d6bEncodeGithubComTheZionMatrixCoreAppAchievementServiceInternalBiz(out *jwriter.Writer, in CommentMap) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"uuid\":"
		out.RawString(prefix[1:])
		out.String(string(in.Uuid))
	}
	{
		const prefix string = ",\"medal\":"
		out.RawString(prefix)
		out.String(string(in.Medal))
	}
	{
		const prefix string = ",\"mode\":"
		out.RawString(prefix)
		out.String(string(in.Mode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CommentMap) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson899f4d6bEncodeGithubComTheZionMatrixCoreAppAchievementServiceInternalBiz(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CommentMap) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson899f4d6bEncodeGithubComTheZionMatrixCoreAppAchievementServiceInternalBiz(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CommentMap) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson899f4d6bDecodeGithubComTheZionMatrixCoreAppAchievementServiceInternalBiz(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CommentMap) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson899f4d6bDecodeGithubComTheZionMatrixCoreAppAchievementServiceInternalBiz(l, v)
}