{{if .CurrentBook -}}
**What are we reading?**

This week we are reading *{{.CurrentBook}}* by {{.CurrentAuthor}}
{{- end}}
{{if .NextBook -}}
We are starting {{.NextBook}} by {{.NextAuthor}} on {{.NextBookStartDate}}
{{- end}}

## 🗓️ Upcoming Schedule

{{range .Schedule -}}
### {{.Date}} ☕️ Meet Up

- **📖 Book**: *{{.BookName}}*
- **📍 Meeting Location**: {{.CafeName}} ([Directions]({{.Link}}))

{{end}}