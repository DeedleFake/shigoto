// shigoto is a very simple static site generator.
//
// shigoto has several basic philosophies:
//
//    - No config files. All configuration is done via directory
//      hierarchy and file metadata.
//    - Content produces output. No automatic categories, no tags,
//      nothing.
//
// Directory Structure
//
// A shigoto project takes a specific structure, directory wise:
//
//    <project root>
//           ┣━━━━━━━ tmpl
//           ┣━━━━━━━ draft
//           ┗━━━━━━━ publish
//
// The tmpl directory stores template information. These describe how
// content is turned into output.
//
// The draft directory stores drafts of content. These will be skipped
// when building a site, but can be published using the "publish"
// command. All files in this directory, regardless of their location
// in subdirectories, are treated the same as if they were in the
// top-level of the directory.
//
// The publish directory stores published content. These will be
// converted into output when the "build" command is run using the
// templates in the tmpl directory. All files in this directory,
// regardless of their location in subdirectories, are treated the
// same as if they were in the top-level of the directory.
//
// Along with these, an option static directory may be included in the
// project root. If this directory exists, any files in it will be
// copied into the output directory verbatim before the actual build
// begins.
//
// File Structure
//
// All files follow a similar structure to each other. Each begins
// with an optional section with metadata, specified as YAML, followed
// by a section with the file's content. These two sections are
// separated by a single line containing nothing except for a minimum
// of five plus signs (+). If the separator line doesn't exist
// anywhere in the file, it is assumed that the file has no metadata
// and the entire file is considered to be content.
//
// A file may contain any metadata that it wants to. The data
// specified is available inside the templates that are used to render
// the file. Some files, however, have several special metadata fields
// that have extra effects on how they are rendered.
//
// In templates, the following fields have an affect:
//
//    - inherit (string): This field specifies which metadata file
//      this one inherits from. If this field is present, then the
//      file that contains the field is rendered and presented to the
//      file that it inherits from in the Content field. The output
//      from that template execution is used as the output of the
//      entire execution. This allows a project to have a global
//      template that provides the basic structure for the site with
//      individual templates that handle specifics.
//
//    - sourceName (string): This field specifies the format to use
//      for creating draft filenames using this template. The contents
//      of this field are themselves executed as a template. The
//      default value is "{{.Title | slug}}.md".
//
//    - buildPath (string): This field specifies the format for the
//      path to an output file in the output directory. The default
//      value is "{{.Title | slug}}/index.{{.Type | ext}}".
//
//    - range ({start: int, end: int, step: int}): This field
//      instructs shigoto to produce a range of files from this
//      template type. This is essentially the same as
//      "for i := start; i < end; i += step { /* Produce file i. */ }".
//      The default values are "{start: 0, end: 1, step: 1}". Note
//      that the range for a template must produce at least one file
//      for the template to do anything.
//
// In drafts, the following fields have an effect:
//
//    - type (string): This field specifies the template type of the
//      draft. There is no default. This field is required.
//
//    - title (string): This field specifies the draft's title. This
//      field is technically optional, but highly recommended. It is
//      used in a number of places, including the creation of file
//      names.
//
// Along with these, any of the fields specified above for templateu
// files can be overriden inside of draft files with the exception of
// "inherit".
//
// Template Execution
//
// When a template is executed, it is passed a data set with any known
// data at the time of its execution. What this means is that every
// template that is executed has access to all of the metadata for any
// files involved as well as several extra pieces of data. For the
// most part, the structure is
//
//    - Type (string): The type of the template involved in this
//      execution.
//
//    - Title (string): The title of the content involved. This can be
//      an empty string in some cases.
//
//    - Tmpl (map): The metadata of the template involved in this
//      execution.
//
//    - Meta (map): The metadata of the content involved in this
//      execution.
//
//    - Range (map): Contains range information. Keys are "Start",
//      "End", and "Current".
//
// Along with these, several functions are available:
//
//    - markdown (string -> string): Runs its input through a Markdown
//      engine and returns the output.
//
//    - slug (string -> string): Converts a string into a slug to make
//      it more suitable for a URL or filename.
//
//    - time (string | int -> time.Time): Parses a time into a
//      time.Time. If it is given an int, it is assumed that that is
//      the number of seconds since the Unix epoch. If it is given a
//      string, parsing of the string is attempted using each of the
//      format constants defined in the Go time package, with the
//      exception of time.Kitchen, in the order that they are
//      specified in that package. If any succeeds then the result is
//      returned.
//
//    - trimExt (string -> string): Trims the extension off of a
//      filename.
//
//    - ext (string -> string): Returns the extension of a filename
//      without the intervening dot.
//
//    - tmpl (string, any -> string): Finds and executes the specified
//      template from the tmpl directory using the given data.
//
//    - pages (string, int -> int): Returns the number of pages
//      required to display all of the content of the given type with
//      the given number of them per page.
package main
