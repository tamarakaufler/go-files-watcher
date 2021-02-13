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

func TestDaemon_IsExcluded(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	type fields struct {
		BasePath  string
		Extention string
		Command   string
		Excluded  []string
		Frequency int32
	}
	type args struct {
		path string
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "file is excluded - regex files exclusion",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/*"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir1/test1.go",
				name: "test1.go",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "file is excluded - regex files exclusion 2",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"test2.go", "fixtures/basepath/subdir1/test1*"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir1/test1.go",
				name: "test1.go",
			},
			want:    true,
			wantErr: false,
		},
		// this does not work - why?
		// {
		// 	name: "file is excluded - regex files exclusion - 3",
		// 	fields: fields{
		// 		BasePath:  "fixtures/basepath",
		// 		Extention: ".go",
		// 		Command:   "echo \"Hello world\"",
		// 		Excluded:  []string{"test2.go", "fixtures/basepath/*/test.go"},
		// 		Frequency: 3,
		// 	},
		// 	args: args{
		// 		path: "fixtures/basepath/subdir1/test.go",
		// 		name: "test.go",
		// 	},
		// 	want:    true,
		// 	wantErr: false,
		// },
		{
			name: "file is excluded - string path exclusion 1",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/test.go"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir1/test.go",
				name: "test.go",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "file is excluded - string file exclusion 2",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"test2.go", "fixtures/basepath/subdir1/test.go"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir2/test2.go",
				name: "test2.go",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "file is excluded - string file exclusion 3",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"test2.go", "fixtures/basepath/subdir1/test.go"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir2/test2.go",
				name: "test2.go",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "file is excluded - regex ? file exclusion",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/test.g?"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir1/test.go",
				name: "test.go",
			},
			want:    true,
			wantErr: false,
		},
		// this does not work - why?
		// {
		// 	name: "file is excluded - regex ? file exclusion 2",
		// 	fields: fields{
		// 		BasePath:  "fixtures/basepath",
		// 		Extention: ".go",
		// 		Command:   "echo \"Hello world\"",
		// 		Excluded:  []string{"fixtures/basepath/subdir1/test.?o"},
		// 		Frequency: 3,
		// 	},
		// 	args: args{
		// 		path: "fixtures/basepath/subdir1/test.go",
		// 		name: "test.go",
		// 	},
		// 	want:    true,
		// 	wantErr: false,
		// },
		{
			name: "file is not excluded - string file exclusion 4",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/test.go"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir1/test2.go",
				name: "test2.go",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "file is not excluded - regex file exclusion 1",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/test*"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir1/aaa.go",
				name: "aaa.go",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "file is not excluded - regex file exclusion 2",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures/basepath/subdir1/*"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir2/aaa.go",
				name: "aaa.go",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "file to be excluded - regex file exclusion",
			fields: fields{
				BasePath:  "fixtures/basepath",
				Extention: ".go",
				Command:   "echo \"Hello world\"",
				Excluded:  []string{"fixtures(a-]basepath/subdir1/*"},
				Frequency: 3,
			},
			args: args{
				path: "fixtures/basepath/subdir2/aaa.go",
				name: "aaa.go",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := daemon.New(
				daemon.WithBasePath(tt.fields.BasePath),
				daemon.WithCommand(tt.fields.Command),
				daemon.WithExcluded(tt.fields.Excluded),
				daemon.WithFrequency(tt.fields.Frequency),
			)
			got, err := d.IsExcluded(ctx, tt.args.path, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.IsExcluded() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Daemon.IsExcluded() = %v, want %v", got, tt.want)
			}
		})
	}
}
