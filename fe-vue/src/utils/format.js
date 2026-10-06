export const formatRupiah = (angka) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(angka)
}

export const formatDate = (dateString) => {
  if (!dateString) return ''
  const [year, month, day] = dateString.split('-')
  if (!year || !month || !day) return dateString
  return `${day}-${month}-${year}`
}

export const toDigits = (value) => String(value ?? '').replace(/\D/g, '')

export const formatThousands = (value) => {
  const digits = toDigits(value)
  return digits ? digits.replace(/\B(?=(\d{3})+(?!\d))/g, '.') : ''
}
