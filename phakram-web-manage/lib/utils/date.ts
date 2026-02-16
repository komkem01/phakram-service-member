export function formatDateOnly(value: string | Date | null | undefined): string {
  if (!value) return '-';

  if (typeof value === 'string') {
    const [datePart] = value.split('T');
    return datePart || '-';
  }

  if (value instanceof Date && !Number.isNaN(value.getTime())) {
    return value.toISOString().split('T')[0];
  }

  return '-';
}
