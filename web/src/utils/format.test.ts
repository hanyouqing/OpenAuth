import { describe, it, expect } from 'vitest'
import { formatDate, formatDateTime, formatFileSize, truncate } from './format'

describe('format', () => {
  describe('formatDate', () => {
    it('should format a date string', () => {
      const date = '2024-01-01T00:00:00Z'
      const formatted = formatDate(date)
      expect(formatted).toContain('2024')
    })

    it('should format a Date object', () => {
      const date = new Date('2024-01-01')
      const formatted = formatDate(date)
      expect(formatted).toContain('2024')
    })
  })

  describe('formatDateTime', () => {
    it('should format a datetime string', () => {
      const date = '2024-01-01T12:00:00Z'
      const formatted = formatDateTime(date)
      expect(formatted).toContain('2024')
    })
  })

  describe('formatFileSize', () => {
    it('should format bytes', () => {
      expect(formatFileSize(0)).toBe('0 Bytes')
      expect(formatFileSize(1024)).toBe('1 KB')
      expect(formatFileSize(1048576)).toBe('1 MB')
    })
  })

  describe('truncate', () => {
    it('should truncate long strings', () => {
      const longString = 'a'.repeat(100)
      const truncated = truncate(longString, 50)
      expect(truncated.length).toBeLessThanOrEqual(53) // 50 + '...'
      expect(truncated).toContain('...')
    })

    it('should not truncate short strings', () => {
      const shortString = 'hello'
      const result = truncate(shortString, 10)
      expect(result).toBe(shortString)
    })
  })
})
