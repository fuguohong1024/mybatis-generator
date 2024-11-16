package {{.Packages}};

{{.Imports}}

@Mapper
public interface {{.JavaBeanName}} {
{{range $i,$v := .Methods}}
    {{$v}}
{{end}}
}