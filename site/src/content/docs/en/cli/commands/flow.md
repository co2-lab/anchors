---
title: "anchors flow & queue"
description: "Operating background watchers, agent task queues, and stage deliveries."
---

Anchors enables continuous, autonomous AI programming through a reactive task queue.

---

## 1. The Watcher & Queue Cycle

1. **`anchors watch`**: Runs in the background. Whenever you or an AI agent saves a file (e.g. `order.spec.md`), the watcher analyzes what's missing and enqueues the next action (`create-feature`, `implement-code`, `run-tests`).
2. **`anchors queue`**: Lists the current tasks waiting in line.
3. **`anchors next`**: The AI agent calls this to claim the next card without needing human prompts.
4. **`anchors work <target>`**: Emits the structured prompt instructing the model what to do for that specific card.
5. **`anchors deliver`**: The model marks completion; Anchors validates the delivery and advances the task.
6. **`anchors done <id>`**: Closes the card and moves it to history.

---

## 2. Handling Obstacles & Escalations

- **`anchors escalate`**: When an AI finds an ambiguous requirement, it parks the card in `needs-user` and opens an issue.
- **`anchors unblock`**: Once the human answers the question in the issue, this unblocks the card.
- **`anchors decided <card>`**: Resumes execution with the human's decision registered.
