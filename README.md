# go-libs

cdcloud-io go libraries

## create a new lib

```bash
cd lib1
go mod init github.com/cdcloud-io/go-libs/lib1
cd ../lib2
go mod init github.com/cdcloud-io/go-libs/lib2
```

## tag the repo

```bash
git tag v1.0.0
git push origin v1.0.0
```

---

### Steps to Individually Version Go Modules in a Single Repository

1. **Create Separate Go Modules for Each Library**:
   Each library should reside in its own directory within the repository, and each directory should have its own `go.mod` file. This allows Go to treat each library as a separate module.

   For example, your repository structure might look like this:

   ```sh
   /my-repo
   ├── /lib1
   │   ├── go.mod
   │   ├── lib1.go
   ├── /lib2
   │   ├── go.mod
   │   ├── lib2.go
   └── /lib3
       ├── go.mod
       ├── lib3.go
   ```

   - `lib1/`, `lib2/`, and `lib3/` are separate libraries, each with its own `go.mod`.
   - Each library can have its own version and dependencies.

2. **Initialize Each Go Module**:
   Navigate to each library directory and initialize a Go module for it:

   ```bash
   cd lib1
   go mod init github.com/cdcloud-io/go-libs/lib1

   cd ../lib2
   go mod init github.com/cdcloud-io/go-libs/lib2

   cd ../lib3
   go mod init github.com/cdcloud-io/go-libs/lib3
   ```

   This will create a `go.mod` file for each library.

3. **Versioning Each Library Individually**:
   In a monorepo, you can tag individual directories with versions using **Git tags**. Go Modules will respect these version tags based on the directory path.

   To version each library, follow these steps:

   - For **lib1**:

     ```bash
     git tag lib1/v1.0.0
     git push origin lib1/v1.0.0
     ```

   - For **lib2**:

     ```bash
     git tag lib2/v1.0.0
     git push origin lib2/v1.0.0
     ```

   - For **lib3**:

     ```bash
     git tag lib3/v1.0.0
     git push origin lib3/v1.0.0
     ```

   Each library will now be versioned separately, and consumers of these libraries can reference them using the corresponding version tag.

4. **Importing and Using Versioned Libraries**:
   Other Go modules can now import your versioned libraries by specifying the path and version.

   For example, to use `lib1` version `v1.0.0`, a consumer's `go.mod` file would include:

   ```go
   require github.com/cdcloud-io/lib1 v1.0.0
   ```

5. **Updating Versions**:
   When you make changes to a specific library and want to release a new version, tag the new version for that specific library.

   - For example, after updating `lib1`:

     ```bash
     git tag lib1/v1.1.0
     git push origin lib1/v1.1.0
     ```

   This will create a new version `v1.1.0` of `lib1`, while other libraries (`lib2`, `lib3`) will remain at their current versions.

### Best Practices

- **Use Semantic Versioning (SemVer)**: Follow semantic versioning for each individual library. For example:
  - `v1.0.0` for an initial release.
  - `v1.1.0` for adding features.
  - `v1.0.1` for bug fixes.
  
- **Maintain Separate `go.mod` Files**: Ensure that each library has its own `go.mod` file and specify dependencies only relevant to that library.

### Example Workflow for New Versions

1. Make changes to `lib1`.
2. Commit the changes.
3. Tag the new version: `git tag lib1/v1.1.0`.
4. Push the tag: `git push origin lib1/v1.1.0`.
5. Consumers can now update their dependencies to use `github.com/cdcloud-io/go-libs/lib1 v1.1.0`.

### TODO

- **CI/CD Automation**: Consider automating the versioning process using a CI/CD pipeline that automatically bumps versions based on changes detected in each library.
