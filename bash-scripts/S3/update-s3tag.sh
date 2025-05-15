#!/bin/bash

# List of bucket names
buckets=("audit-vj-citrix-users" "audit-vj-travelagents")

for bucket in "${buckets[@]}"; do
  echo "Update S3 bucket tags for $bucket"

  # Update the tags
  aws s3api put-bucket-tagging \
    --bucket "$bucket" \
    --tagging file://"${bucket}-tags.json"

  if [ $? -eq 0 ]; then
    echo "Tags updated successfully for $bucket from ${bucket}-tags.json"
  else
    echo "Error updating tags for $bucket"
  fi

  echo "----------------------------------------"
done
