import { ref } from 'vue'
import { useApi } from './useApi'

export function useExamples() {
  const examples = ref([])
  const example = ref(null)
  const { loading, error, get, post, put, del } = useApi()

  const fetchExamples = async () => {
    try {
      const data = await get('/examples')
      examples.value = data.data || []
    } catch (err) {
      console.error('Gagal mengambil data:', err)
    }
  }

  const fetchExample = async (id) => {
    try {
      const data = await get(`/examples/${id}`)
      example.value = data.data
    } catch (err) {
      console.error('Gagal mengambil detail:', err)
    }
  }

  const createExample = async (formData) => {
    try {
      await post('/examples', formData)
      await fetchExamples()
      return true
    } catch (err) {
      console.error('Gagal membuat data:', err)
      return false
    }
  }

  const updateExample = async (id, formData) => {
    try {
      await put(`/examples/${id}`, formData)
      await fetchExamples()
      return true
    } catch (err) {
      console.error('Gagal memperbarui data:', err)
      return false
    }
  }

  const deleteExample = async (id) => {
    try {
      await del(`/examples/${id}`)
      await fetchExamples()
      return true
    } catch (err) {
      console.error('Gagal menghapus data:', err)
      return false
    }
  }

  return {
    examples,
    example,
    loading,
    error,
    fetchExamples,
    fetchExample,
    createExample,
    updateExample,
    deleteExample
  }
}
