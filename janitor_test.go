package cache

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJanitorRun(t *testing.T) {
	cache := New(DefaultExpiration, 50*time.Millisecond)
	cache.Set("key1", "value1", 10*time.Millisecond)
	cache.Set("key2", "value2", NoExpiration)

	time.Sleep(100 * time.Millisecond) // Esperar a que los elementos expiren

	cache.mu.RLock()
	_, found1 := cache.items["key1"]
	_, found2 := cache.items["key2"]
	cache.mu.RUnlock()

	assert.False(t, found1, "key1 debería haber expirado y ser eliminado por el janitor")
	assert.True(t, found2, "key2 no debería haber expirado y debería seguir presente")
}

func TestStopJanitor(t *testing.T) {
	cache := New(DefaultExpiration, 50*time.Millisecond)
	cache.Set("key", "value", 10*time.Millisecond)

	stopJanitor(cache)
	runtime.GC() // Forzar el recolector de basura para ejecutar el finalizador

	time.Sleep(100 * time.Millisecond) // Esperar suficiente tiempo para confirmar que el janitor está detenido

	cache.mu.RLock()
	_, found := cache.items["key"]
	cache.mu.RUnlock()

	// Aquí `found` debe ser `true` porque el janitor se detuvo antes de poder eliminar el elemento expirado
	assert.True(t, found, "key debería seguir presente ya que el janitor fue detenido")
}
