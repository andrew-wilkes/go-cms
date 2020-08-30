<?php
Tokens::add(['#TITLE#' => Page::get_value('title')]);

$posts = Pages::get_recent_posts();

if (count($posts))
	HTML::script_meta_tags(['canonical' => $posts[0]->key]);

include 'header.php';

new RecentPosts();

?>
<div class="sidebar">

<h3>Recent Posts</h3>

#RECENT#

<h3>Archives</h3>

<?php Archive::get_content(); ?>

<h3>Categories</h3>

<?php (new Categories())->print_list(); ?>

</div>

<div class="content right">
<?php

if (count($posts))
{ ?>
	<div>
		<h1><?php echo HTML::hyperlink($posts[0]->key, $posts[0]->title); ?></h1>
		<p class="meta">Published: <?php echo date('jS F Y', strtotime($posts[0]->published)); ?></p>
	</div>

	<div>
		<?php echo Page::insert_region_content('main-content', $posts[0]->id); ?>
		<p><?php echo HTML::hyperlink($posts[0]->key . '#comments', 'Comments'); ?></p>
	</div>
<?php
}
else
	echo "<p>No posts have been published yet.</p>";
?>
</div>

<?php include 'footer.php'; ?>