<template>
  <b-container
    fluid
    data-test-id="editor-toolbar"
    class="bg-white shadow border-top p-3"
  >
    <b-row
      no-gutters
      class="wrap-with-vertical-gutters align-items-center"
    >
      <div
        class="wrap-with-vertical-gutters align-items-center"
      >
        <b-button
          data-test-id="button-back-without-save"
          variant="link"
          :disabled="buttonProcessing"
          class="text-dark back mr-auto"
          @click="$emit('back')"
        >
          <font-awesome-icon
            :icon="['fas', 'chevron-left']"
            class="back-icon"
          />
          {{ $t('label.backWithoutSave') }}
        </b-button>
      </div>
      <div
        class="ml-auto"
      >
        <slot />
      </div>
      <div
        class="d-flex wrap-with-vertical-gutters align-items-center ml-auto"
      >
        <slot name="delete" />
        <c-input-confirm
          v-if="!hideDelete"
          v-b-tooltip.hover
          :disabled="disableDelete || buttonProcessing"
          size="lg"
          size-confirm="lg"
          variant="danger"
          :title="deleteTooltip"
          :borderless="false"
          @confirmed="$emit('delete')"
        >
          {{ $t('label.delete') }}
        </c-input-confirm>
        <slot name="saveAsCopy" />

        <c-button-submit
          v-if="!hideClone"
          data-test-id="button-clone"
          :disabled="disableClone || buttonCloneProcessing"
          :title="cloneTooltip"
          :processing="buttonCloneProcessing"
          :variant="'light'"
          :button-text="$t('label.saveAsCopy')"
          :button-size="'lg'"
          class="ml-2"
          @submit="$emit('clone')"
        />
        <c-button-submit
          v-if="!hideSave"
          data-test-id="button-save-and-close"
          :disabled="disableSave || buttonSaveAndCloseProcessing"
          :processing="buttonSaveAndCloseProcessing"
          :variant="'light'"
          :button-text="$t('label.saveAndClose')"
          :button-size="'lg'"
          class="ml-2"
          @submit="$emit('saveAndClose')"
        />
        <c-button-submit
          v-if="!hideSave"
          data-test-id="button-save"
          :disabled="disableSave || buttonSaveProcessing"
          :processing="buttonSaveProcessing"
          :button-text="$t('label.save')"
          :button-size="'lg'"
          class="ml-2"
          @submit="$emit('save')"
        />
      </div>
    </b-row>
  </b-container>
</template>
<script>

import { components } from '@cortezaproject/corteza-vue'
const { CButtonSubmit } = components

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    CButtonSubmit,
  },

  inheritAttrs: true,

  props: {
    buttonProcessing: {
      type: Boolean,
      default: false,
    },

    buttonSaveProcessing: {
      type: Boolean,
      default: false,
    },

    buttonSaveAndCloseProcessing: {
      type: Boolean,
      default: false,
    },

    buttonCloneProcessing: {
      type: Boolean,
      default: false,
    },

    backLink: {
      type: Object,
      required: false,
      default: undefined,
    },

    hideDelete: {
      type: Boolean,
      required: false,
    },

    hideSave: {
      type: Boolean,
      required: false,
    },

    hideClone: {
      type: Boolean,
      required: false,
    },

    disableDelete: {
      type: Boolean,
      required: false,
      default: false,
    },

    disableSave: {
      type: Boolean,
      required: false,
      default: false,
    },

    disableClone: {
      type: Boolean,
      default: false,
    },

    deleteTooltip: {
      type: String,
      required: false,
      default: '',
    },

    cloneTooltip: {
      type: String,
      default: '',
    },
  },
}
</script>
<style lang="scss" scoped>
.back {
  &:hover {
    text-decoration: none;

    .back-icon {
      transition: transform 0.3s ease-out;
      transform: translateX(-4px);
    }
  }
}

[dir="rtl"] {
  .back {
    .back-icon {
      display: none;
    }
  }
}
</style>
