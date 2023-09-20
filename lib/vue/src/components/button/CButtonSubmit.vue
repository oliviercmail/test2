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
        v-if="loadingText"
        class="loading-text mx-2"
      >
        {{ loadingText }}
      </span>
      <b-spinner
        v-else
        small
        :variant="iconVariant"
      />
    </template>
    <template v-else-if="success">
      <font-awesome-icon
        :class="iconVariant"
        class="text-white h3 mb-0"
      />
    </template>
    <template v-else>
      <span>
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
