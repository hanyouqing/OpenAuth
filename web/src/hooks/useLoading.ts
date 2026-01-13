import { useState, useCallback } from 'react'

export function useLoading() {
  const [loading, setLoading] = useState(false)

  const withLoading = useCallback(
    async <T,>(asyncFn: () => Promise<T>): Promise<T | undefined> => {
      setLoading(true)
      try {
        return await asyncFn()
      } finally {
        setLoading(false)
      }
    },
    []
  )

  return { loading, withLoading }
}
