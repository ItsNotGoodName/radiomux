import { Accessor, batch, createMemo, createSignal } from "solid-js"

export type UseSelectionResult = {
  record: Accessor<Record<number, boolean>>
  ids: Accessor<Array<number>>
  isMultiple: Accessor<boolean>
  check: (id: number, selected: boolean) => void
  checkMultiple: (selected: boolean) => void
  clear: () => void
}

export function useSelection(data: Accessor<Array<{ id: number }>>): UseSelectionResult {
  const [_record, setRecord] = createSignal<Record<number, boolean>>({})
  const [_isMultiple, setIsMultiple] = createSignal(false)
  const memo = createMemo((): [Record<number, boolean>, Array<number>, boolean] => {
    let record: Record<number, boolean> = {}
    let ids: Array<number> = []
    let isMultiple = data().length > 0
    for (const d of data()) {
      if (_record()[d.id] == true) {
        record[d.id] = true
        ids.push(d.id)
      } else {
        isMultiple = false
        record[d.id] = false
      }
    }
    return [record, ids, isMultiple]
  })
  const record = () => memo()[0]
  const ids = () => memo()[1]
  const isMultiple = () => memo()[2]

  const check = (id: number, checked: boolean) => {
    let value: Record<number, boolean> = {}
    for (const d of data()) {
      value[d.id] = record()[d.id] == true
    }
    value[id] = checked
    setRecord(value)
  }

  const checkMultiple = (checked: boolean) => {
    batch(() => {
      let value: Record<number, boolean> = {}
      for (const d of data()) {
        value[d.id] = checked
      }
      setIsMultiple(checked)
      setRecord(value)
    })
  }

  const clear = () => {
    batch(() => {
      setIsMultiple(false)
      setRecord({})
    })
  }

  return {
    record,
    ids,
    isMultiple,
    check,
    checkMultiple,
    clear,
  }
}

