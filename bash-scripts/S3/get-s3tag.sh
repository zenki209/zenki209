#!/bin/bash

# List of bucket names
buckets=("audit-vj-citrix-users" "audit-vj-travelagents")
for bucket in "${buckets[@]}"; do
  echo "Fetching tags for bucket: $bucket"

  # Get the tag set (if any)
  tags=$(aws s3api get-bucket-tagging --bucket "$bucket" 2>/dev/null)

  if [ $? -eq 0 ]; then
    echo "Saving tags for $bucket to ${bucket}-tags.json"
    echo "$tags" | jq . > "${bucket}-tags.json"
  else
    echo "No tags found or error accessing tags for $bucket"
  fi

  echo "----------------------------------------"
done

