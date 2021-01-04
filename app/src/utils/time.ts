import { format, parseISO } from "date-fns";

const dateFmtString = "yyyy-MM-dd HH:mm:ss OOOO"
export function formatCreated(iso: string) {
  const time = parseISO(iso)
  
  return `Created: ${format(time, dateFmtString)}`
}

export function formatUpdated(iso: string) {
  const time = parseISO(iso)
  
  return `Updated: ${format(time, dateFmtString)}`
}