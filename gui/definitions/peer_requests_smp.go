package definitions

func init() {
	add(`PeerRequestsSMP`, &defPeerRequestsSMP{})
}

type defPeerRequestsSMP struct{}

func (*defPeerRequestsSMP) String() string {
	return `<interface>
  <object class="GtkInfoBar" id="peer_requests_smp">
    <property name="message-type">GTK_MESSAGE_WARNING</property>
    <child internal-child="content_area">
      <object class="GtkBox" id="box">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkLabel" id="message">
            <property name="wrap">true</property>
          </object>
        </child>
      </object>
      <object class="GtkButtonBox" id="button_box">
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkButton" id="verification_button">
            <property name="label" translatable="yes">Finish verification</property>
          </object>
        </child>
        <child>
          <object class="GtkButton" id="cancel_button">
            <property name="label" translatable="yes">Cancel</property>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}
