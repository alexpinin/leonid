#!/bin/bash
set -e

# If arguments passed — use them as whisper args directly
if [ $# -gt 0 ]; then
    exec whisper "$@"
fi

# Otherwise, process all files in /data
INPUT_DIR="/data"
OUTPUT_DIR="/output"

shopt -s nullglob
files=("$INPUT_DIR"/*.{ogg,mp3,wav,m4a,flac,webm,mp4})
shopt -u nullglob

if [ ${#files[@]} -eq 0 ]; then
    echo "No audio files found in /data"
    echo "Mount a directory with audio files: -v /path/to/audio:/data"
    exit 1
fi

echo "Found ${#files[@]} file(s), model=${WHISPER_MODEL}, lang=${WHISPER_LANG}"
echo "---"

for file in "${files[@]}"; do
    filename=$(basename "$file")
    echo "Processing: $filename"
    whisper "$file" \
        --model "$WHISPER_MODEL" \
        --language "$WHISPER_LANG" \
        --output_format "txt" \
        --output_dir "$OUTPUT_DIR"
    echo "Done: $filename"
    echo "---"
done

echo "All files processed. Results in /output"