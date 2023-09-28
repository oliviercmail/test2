<template>
  <b-button
    :data-test-id="cypressID"
    type="submit"
    :variant="variant"
    :disabled="disabled || processing || success"
    :size="size"
    :block="block"
    :title="title"
    :class="buttonClass"
    @click.prevent="$emit('submit')"
  >
    <template v-if="processing">
      <span
        data-test-id="button-loading-text"
        v-if="loadingText"
        class="loading-text mx-2"
      >
        {{ loadingText }}
      </span>
      <b-spinner
        data-test-id="spinner"
        v-else
        small
      />
    </template>
    <template v-else-if="success">
      <font-awesome-icon
        data-test-id="icon-success"
        :icon="['fas', 'check']"
        :class="iconVariant"
        class="text-white"
      />
    </template>
    <template v-else>
      <span
        data-test-id="button-text"
      >
        {{ text }}
      </span>
    </template>
  </b-button>
</template>

<script>
export default {
  name: 'CSubmitButton',

  i18nOptions: {
    namespaces: 'admin',
  },

  props: {
    processing: {
      type: Boolean,
      value: false,
    },

    success: {
      type: Boolean,
      value: false,
    },

    disabled: {
      type: Boolean,
      value: false,
    },

    title: {
      type: String,
      default: '',
    },

    buttonClass: {
      type: String,
      default: '',
    },

    text: {
      type: String,
      default: '',
    },

    loadingText: {
      type: String,
      default: '',
    },

    size: {
      type: String,
      default: 'md',
    },

    block: {
      type: Boolean,
      value: false,
    },

    variant: {
      type: String,
      default: 'primary',
    },

    iconVariant: {
      type: String,
      default: 'text-white',
    },

    cypressID: {
      type: String,
      default: 'button-submit',
    },
  },
}
</script>

<style>
.submit {
  min-width: 75px;
  min-height: 35px;
}

.loading-text::after {
  display: inline-block;
  animation: saving steps(1, end) 1s infinite;
  content: '';
}
</style>
