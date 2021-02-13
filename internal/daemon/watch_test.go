// +build unit_test

package daemon_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/tamarakaufler/go-files-watcher/internal/daemon"
)

func TestDaemon_CollectFiles_HappyPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	type fields struct {
		BasePath  string
		Extention string
		Command   string
		Excluded  []string
		Frequency int32
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "got correctly all files - no exclusions",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{},
				Frequency: 3,
			},
			want: []string{"test.go", "test1.go", "test.go", "test2.go"},
		},
		{
			name: "got correctly all files - with one individual file exclusion",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/test.go"},
				Frequency: 3,
			},
			want: []string{"test1.go", "test.go", "test2.go"},
		},
		{
			name: "got correctly all files - with individual file exclusions",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/test.go", "fixtures/basepath/subdir2/test2.go"},
				Frequency: 3,
			},
			want: []string{"test1.go", "test.go"},
		},
		{
			name: "got correctly all files - with regex file exclusions",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/*", "fixtures/basepath/subdir2/test.go"},
				Frequency: 3,
			},
			want: []string{"test2.go"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			d := daemon.New(
				daemon.WithBasePath(tt.fields.BasePath),
				daemon.WithCommand(tt.fields.Command),
				daemon.WithExcluded(tt.fields.Excluded),
				daemon.WithFrequency(tt.fields.Frequency),
			)

			got := d.CollectFiles(ctx)
			gotNames := extractNames(got)
			if !reflect.DeepEqual(gotNames, tt.want) {
				t.Errorf("Daemon.CollectFiles() = %v, want %v", gotNames, tt.want)
			}
		})
	}
}

func extractNames(files []os.FileInfo) []string {
	names := []string{}
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names
}
